package integration_tests

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/kubernetes-sig-testing/frameworks/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The Testing Framework", func() {
	var controlPlane *integration.ControlPlane

	AfterEach(func() {
		Expect(controlPlane.Stop()).To(Succeed())
	})

	It("Successfully manages the control plane lifecycle", func() {
		var err error

		controlPlane = &integration.ControlPlane{}

		By("Starting all the control plane processes")
		err = controlPlane.Start()
		Expect(err).NotTo(HaveOccurred(), "Expected controlPlane to start successfully")

		apiServerURL := controlPlane.APIURL()
		etcdClientURL := controlPlane.APIServer.EtcdURL
		controllerManagerURL := controlPlane.ControllerManager.URL
		schedulerURL := controlPlane.Scheduler.URL

		isEtcdListeningForClients := isSomethingListeningOnPort(etcdClientURL.Host)
		isAPIServerListening := isSomethingListeningOnPort(apiServerURL.Host)
		isControllerManagerListening := isSomethingListeningOnPort(controllerManagerURL.Host)
		isSchedulerListening := isSomethingListeningOnPort(schedulerURL.Host)

		By("Ensuring Etcd is listening")
		Expect(isEtcdListeningForClients()).To(BeTrue(),
			fmt.Sprintf("Expected Etcd to listen for clients on %s,", etcdClientURL.Host))

		By("Ensuring APIServer is listening")
		Expect(isAPIServerListening()).To(BeTrue(),
			fmt.Sprintf("Expected APIServer to listen on %s", apiServerURL.Host))

		By("Ensuring ControllerManager is listening")
		Expect(isControllerManagerListening()).To(BeTrue(),
			fmt.Sprintf("Expected ControllerManager to listen on %s", controllerManagerURL.Host))

		By("Ensuring Scheduler is listening")
		Expect(isSchedulerListening()).To(BeTrue(),
			fmt.Sprintf("Expected Scheduler to listen on %s", schedulerURL.Host))

		By("Starting a virtual kubelet")
		vk := &integration.VirtualKubelet{}
		vk.APIServerURL = apiServerURL
		Expect(vk.Start()).NotTo(HaveOccurred())
		defer func() {
			By("Stopping the virtual kubelet")
			Expect(vk.Stop()).To(Succeed())
		}()

		By("Getting a kubectl & run it against the control plane")
		kubeCtl := controlPlane.KubeCtl()
		stdout, stderr, err := kubeCtl.Run("get", "nodes")
		Expect(err).NotTo(HaveOccurred())
		errBytes, err := ioutil.ReadAll(stderr)
		Expect(err).NotTo(HaveOccurred())
		Expect(stdout).To(ContainSubstring("virtual-kubelet   Ready"))
		Expect(errBytes).To(BeEmpty())

		By("Stopping all the control plane processes")
		err = controlPlane.Stop()
		Expect(err).NotTo(HaveOccurred(), "Expected controlPlane to stop successfully")

		By("Ensuring Etcd is not listening anymore")
		Expect(isEtcdListeningForClients()).To(BeFalse(), "Expected Etcd not to listen for clients anymore")

		By("Ensuring APIServer is not listening anymore")
		Expect(isAPIServerListening()).To(BeFalse(), "Expected APIServer not to listen anymore")

		By("Ensuring ControllerManager is not listening anymore")
		Expect(isControllerManagerListening()).To(BeFalse(), "Expected ControllerManager not to listen anymore")

		By("Ensuring Scheduler is not listening anymore")
		Expect(isSchedulerListening()).To(BeFalse(), "Expected Scheduler not to listen anymore")

		By("Not erroring when stopping a stopped ControlPlane")
		Expect(func() {
			Expect(controlPlane.Stop()).To(Succeed())
		}).NotTo(Panic())
	})

	Context("when Stop() is called on the control plane", func() {
		Context("but the control plane is not started yet", func() {
			It("does not error", func() {
				controlPlane = &integration.ControlPlane{}

				stoppingTheControlPlane := func() {
					Expect(controlPlane.Stop()).To(Succeed())
				}

				Expect(stoppingTheControlPlane).NotTo(Panic())
			})
		})
	})

	Context("when the control plane is configured with its components", func() {
		It("it does not default them", func() {
			myEtcd, myAPIServer :=
				&integration.Etcd{StartTimeout: 15 * time.Second},
				&integration.APIServer{StopTimeout: 16 * time.Second}

			controlPlane = &integration.ControlPlane{
				Etcd:      myEtcd,
				APIServer: myAPIServer,
			}

			Expect(controlPlane.Start()).To(Succeed())
			Expect(controlPlane.Etcd).To(BeIdenticalTo(myEtcd))
			Expect(controlPlane.APIServer).To(BeIdenticalTo(myAPIServer))
			Expect(controlPlane.Etcd.StartTimeout).To(Equal(15 * time.Second))
			Expect(controlPlane.APIServer.StopTimeout).To(Equal(16 * time.Second))
		})
	})

	Measure("It should be fast to bring up and tear down the control plane", func(b Benchmarker) {
		b.Time("lifecycle", func() {
			controlPlane = &integration.ControlPlane{}

			controlPlane.Start()
			controlPlane.Stop()
		})
	}, 10)
})

type portChecker func() bool

func isSomethingListeningOnPort(hostAndPort string) portChecker {
	return func() bool {
		conn, err := net.DialTimeout("tcp", hostAndPort, 1*time.Second)

		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
}
