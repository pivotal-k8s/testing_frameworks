#!/usr/bin/env python

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import boilerplate
import unittest
import StringIO
import os
import sys
import tempfile

base_dir = os.getcwd()

class DefaultArgs(object):
  def __init__(self):
    self.filenames = []
    self.rootdir = "."
    self.boilerplate_dir = base_dir
    self.verbose = True

class TestBoilerplate(unittest.TestCase):
  """
  Note: run this test from the hack/boilerplate directory.

  $ python -m unittest boilerplate_test
  """
  def setUp(self):
    os.chdir(base_dir)
    boilerplate.args = DefaultArgs()

  def test_boilerplate(self):
    os.chdir("test/")

    # capture stdout
    old_stdout = sys.stdout
    sys.stdout = StringIO.StringIO()

    ret = boilerplate.main()

    output = sorted(sys.stdout.getvalue().split())

    sys.stdout = old_stdout

    self.assertEquals(
        output, ['././fail.go', '././fail.py'])

  def test_read_config(self):
    config_file = "./test_with_config_file/boilerplate.json"
    config = boilerplate.read_config_file(config_file)
    self.assertEqual(config.get('dirs_to_skip'), ['dir_to_skip', 'dont_want_this', 'not_interested', '.'])
    self.assertEqual(config.get('not_generated_files_to_skip'), ['alice skips a file', 'bob skips another file'])

  def test_read_nonexistent_config(self):
    config_file = '/nonexistent'
    config = boilerplate.read_config_file(config_file)
    self.assertEqual(config['dirs_to_skip'], boilerplate.default_skipped_dirs)
    self.assertEqual(config['not_generated_files_to_skip'], boilerplate.default_skipped_not_generated)

  def test_read_malformed_config(self):
    config_file = './test_with_config_file/boilerplate.bad.json'
    with self.assertRaises(Exception):
      boilerplate.read_config_file(config_file)

  def test_read_config_called_with_correct_path(self):
    self.has_been_called = False

    def fake_read_config_file(config_file_path):
      self.assertEqual(config_file_path, "/tmp/some/path/boilerplate.json")
      self.has_been_called = True
      return {}

    def nonParallelSafeSetUp():
      self.real_read_config_file = boilerplate.read_config_file
      boilerplate.read_config_file = fake_read_config_file

    def nonParallelSafeTearDown():
      boilerplate.read_config_file = self.real_read_config_file

    nonParallelSafeSetUp()

    try:
      boilerplate.args.rootdir = "/tmp/some/path"
      boilerplate.main()
      self.assertEqual(self.has_been_called, True)
    finally:
      nonParallelSafeTearDown()

  def test_get_files_with_skipping_dirs(self):
    refs = boilerplate.get_refs()
    skip_dirs = ['.']
    files = boilerplate.get_files(refs, skip_dirs)

    self.assertEqual(files, [])

  def test_get_files_with_skipping_not_generated_files(self):
    refs = boilerplate.get_refs()
    regexes = boilerplate.get_regexs()
    files_to_skip = ['boilerplate.py']
    filename = 'boilerplate.py'

    passes = boilerplate.file_passes(filename, refs, regexes, files_to_skip)

    self.assertEqual(passes, True)

  def test_ignore_when_no_valid_boilerplate_template(self):
    with tempfile.NamedTemporaryFile() as temp_file_to_check:
      passes = boilerplate.file_passes(temp_file_to_check.name, boilerplate.get_refs(), boilerplate.get_regexs(), [])
      self.assertEqual(passes, True)

