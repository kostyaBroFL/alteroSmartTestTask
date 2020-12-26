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
go_register_toolchains(version = "1.15.6")

# LICENSE: Apache 2.0 (https://github.com/bazelbuild/rules_proto/blob/master/LICENSE).
http_archive(
    name = "rules_proto",
    sha256 = "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
    strip_prefix = "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
        "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
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
