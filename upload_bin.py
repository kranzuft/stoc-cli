import sys
import zipfile
import os

github_rest_command = """
curl
  -X POST
  -H "Accept: application/vnd.github+json"
  -H "Authorization: token <github_token>"
  https://api.github.com/repos/kranzuft/releases
  -d '{"tag_name":"v<release_version>","target_commitish":"main","name":"v<release_version>","body":"<release_description>","draft":false,"prerelease":false,"generate_release_notes":false}'
""".replace("\n", " ")


def is_blank(str_arg):
    return str_arg == '' or len(str.strip(str_arg, ' ')) == 0


def prepare_builds(build_folder, build_targets):
    for build in build_targets:
        target = build_targets[build]
        with zipfile.ZipFile(buildFolder + "/" + target, mode="w") as archive:
            archive.write(build_folder + "/" + build, os.path.basename(build))


def prepare_api_call(args, call):
    github_token = args[1]
    release_version = args[2]
    release_description = args[3]
    call.replace("<release_version>", release_version) \
        .replace("<release_description>", release_description) \
        .replace("<github_token>", github_token)

    return call


if len(sys.argv) != 4 or is_blank(sys.argv[1]) or is_blank(sys.argv[2]) or is_blank(sys.argv[3]):
    print("Need 3 non-blank arguments")
else:
    if not os.path.isdir('bin/zips'):
        os.mkdir('bin/zips')
    builds = {
        'darwin/amd64/go_build_darwin_64': 'zips/soct_mac_64.zip',
        'linux/386/go_build_linux_32_linux': 'zips/soct_linux_32.zip',
        'linux/amd64/go_build_linux_64_linux': 'zips/soct_linux_64.zip',
        'windows/386/go_build_windows_32.exe': 'zips/soct_windows_32.zip',
        'windows/amd64/go_build_windows_64.exe': 'zips/soct_windows_64.zip'
    }

    buildFolder = "bin"

    prepare_builds(buildFolder, builds)
