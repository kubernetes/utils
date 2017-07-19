# Utils

[![Build Status](https://travis-ci.org/kubernetes/utils.svg?branch=master)](https://travis-ci.org/kubernetes/utils)

Utils is a set of golang libraries that are not specific to
Kubernetes. They should be available and useful to any other Go project
out there.

## Purpose

The Kubernetes project uses a lot of Golang patterns that are re-used in many
difference places, and as the project is further split into different
repositories, this repository is the perfect place to hold this common code.

## Goals

Go libraries in this repository must be:

- Generic enough that they are useful for external/non-kubernetes
  projects,
- Well factored, well tested and reliable,
- Be completely go compliant (go get/build/test/etc)
- Have enough complexity to be shared,
- Have stable APIs, or backward compatible.

The goal is to keep libraries organized in logical entities.

## Libraries

- [Exec](/exec) provides an interface for `os/exec`. It makes it easier
  to mock and replace in tests, especially with
  the [FakeExec](exec/testing/fake_exec.go) struct.
