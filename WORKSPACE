workspace(name = "alteroSmartTestTask")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# LICENSE: Apache 2.0 (https://github.com/bazelbuild/rules_go/blob/master/LICENSE.txt).
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7904dbecbaffd068651916dce77ff3437679f9d20e1a7956bff43826e7645fcc",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
    ],
)
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()
go_register_toolchains(go_version = "1.15.6")

http_archive(
    name = "rules_proto",
    sha256 = "7994a4587e00b9049fed87390fc0d5ff62e0077c1ae8a0761d618e4dce2c525c",
    strip_prefix = "rules_proto-33549b80b8097502de2a966d764c8d23c59f4d08",
    urls = [
        # master, as of 2019-11-04
        "https://github.com/bazelbuild/rules_proto/archive/33549b80b8097502de2a966d764c8d23c59f4d08.tar.gz",
    ],
)
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
rules_proto_dependencies()
rules_proto_toolchains()

http_archive(
    name = "bazel_gazelle",
    sha256 = "222e49f034ca7a1d1231422cdb67066b885819885c356673cb1f72f748a3c9d4",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
    ],
)
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

http_archive(
    name = "com_github_mwitkow_go_proto_validators",
    patch_args = ["-p1"],
    patches = ["//third_party:com_github_mwitkow_go_proto_validators/com_github_mwitkow_go_proto_validators-gazelle.patch"],
    sha256 = "0b36dc401558728fb5a796e48be3251078aabe2f759a56bb2337476873232fc5",
    strip_prefix = "go-proto-validators-fbdcedf3a5550890154208a722600dd6af252902",
    urls = [
        "https://github.com/mwitkow/go-proto-validators/archive/fbdcedf3a5550890154208a722600dd6af252902.zip",  # master, as of 2019-03-03
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains")

load("@io_bazel_rules_go//proto:gogo.bzl", "gogo_special_proto")

http_archive(
    name = "org_golang_google_grpc",
    patch_args = ["-p1"],
    patches = [
        "//third_party:org_golang_google_grpc/org_golang_google_grpc-gazelle.patch",
    ],
    sha256 = "99c2d8e0392b938ab9867a60451964132a43ef67b2a15f8273544145be8006ff",
    strip_prefix = "grpc-go-1.24.0",
    type = "zip",
    urls = ["https://github.com/grpc/grpc-go/archive/v1.24.0.zip"],
    # gazelle args: -go_prefix google.golang.org/grpc
)

http_archive(
    name = "org_golang_x_net",
    patch_args = ["-p1"],
    patches = [
        "//third_party:org_golang_x_net/org_golang_x_net-gazelle.patch",
    ],
    sha256 = "e3ecef768e2834542bdefaf99d047b18bfc81c1ffe5195571102f7ce863a2edf",
    strip_prefix = "net-0deb6923b6d97481cb43bc1043fe5b72a0143032",
    type = "zip",
    # master, as of 2019-11-01
    urls = ["https://github.com/golang/net/archive/0deb6923b6d97481cb43bc1043fe5b72a0143032.zip"],
    # gazelle args: -go_prefix golang.org/x/net
)
