#!/usr/bin/env bash

crane push ./test/e2e/testdata/wolfi-base-docker.tar localhost:5000/wolfi-base:latest
crane push ./test/e2e/testdata/wolfi-base-att/ localhost:5000/wolfi-base:sha256-4d31ef1460be2813657ce7ab3cfd0df2a7366a9b72732d4978b2794cbeb8cd32.att
crane push ./test/e2e/testdata/wolfi-base-sig/ localhost:5000/wolfi-base:sha256-4d31ef1460be2813657ce7ab3cfd0df2a7366a9b72732d4978b2794cbeb8cd32.sig
crane push ./test/e2e/testdata/wolfi-base-docker.tar localhost:5000/notsigned:latest
crane push ./test/e2e/testdata/alpine-cves localhost:5000/alpine-cves
crane push ./test/e2e/testdata/alpine-cves.att localhost:5000/alpine-cves:sha256-eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978.att

crane digest localhost:5000/wolfi-base:latest
crane digest localhost:5000/wolfi-base:sha256-4d31ef1460be2813657ce7ab3cfd0df2a7366a9b72732d4978b2794cbeb8cd32.att
crane digest localhost:5000/wolfi-base:sha256-4d31ef1460be2813657ce7ab3cfd0df2a7366a9b72732d4978b2794cbeb8cd32.sig
crane digest localhost:5000/alpine-cves
crane digest localhost:5000/alpine-cves:sha256-eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978.att

cosign tree localhost:5000/wolfi-base:latest
cosign tree localhost:5000/alpine-cves
cosign verify localhost:5000/wolfi-base@sha256:4d31ef1460be2813657ce7ab3cfd0df2a7366a9b72732d4978b2794cbeb8cd32 --certificate-identity=puerco@chainguard.dev --certificate-oidc-issuer=https://accounts.google.com