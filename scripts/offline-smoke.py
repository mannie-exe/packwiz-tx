"""Exercise the fork's offline batch contract against a built binary."""
import hashlib
from pathlib import Path
import subprocess
import sys
import tempfile
import tomllib

binary = str(Path(sys.argv[1]).resolve())
with tempfile.TemporaryDirectory(prefix="packwiz-tx-smoke-") as temp:
    root = Path(temp)
    index = b'hash-format = "sha256"\n'
    (root / "index.toml").write_bytes(index)
    (root / "pack.toml").write_text(
        'name = "Offline fixture"\npack-format = "packwiz:1.1.0"\n'
        '[index]\nfile = "index.toml"\nhash-format = "sha256"\n'
        f'hash = "{hashlib.sha256(index).hexdigest()}"\n'
        '[versions]\nminecraft = "1.21.1"\n'
    )
    pack = (root / "pack.toml").read_bytes()

    def run(*args, success=True):
        result = subprocess.run(
            [binary, "--config", str(root / "config.toml"), "--cache", str(root / "cache"), *args],
            cwd=root, capture_output=True, text=True, timeout=30,
        )
        assert (result.returncode == 0) == success, result.stdout + result.stderr

    common = ["--offline", "--no-refresh", "--hash", "a" * 64, "--hash-format", "sha256"]
    run("modrinth", "add", *common, "--project-id", "Project1", "--version-id", "Version1",
        "--name", "Fixture MR", "--slug", "fixture-mr", "--filename", "mr.jar",
        "--url", "https://example.invalid/mr.jar", "--side", "client")
    run("curseforge", "add", *common, "--addon-id", "123", "--file-id", "456",
        "--name", "Fixture CF", "--slug", "fixture-cf", "--filename", "cf.jar")
    assert (root / "index.toml").read_bytes() == index
    assert (root / "pack.toml").read_bytes() == pack
    mr = tomllib.loads((root / "mods/fixture-mr.pw.toml").read_text())
    cf = tomllib.loads((root / "mods/fixture-cf.pw.toml").read_text())
    assert mr["side"] == "client" and mr["update"]["modrinth"]["mod-id"] == "Project1"
    assert mr["update"]["modrinth"]["version"] == "Version1"
    assert cf["update"]["curseforge"] == {"project-id": 123, "file-id": 456}
    assert cf["download"]["mode"] == "metadata:curseforge"
    run("refresh")
    entries = tomllib.loads((root / "index.toml").read_text())["files"]
    assert {e["file"] for e in entries} == {"mods/fixture-mr.pw.toml", "mods/fixture-cf.pw.toml"}
    before = (root / "index.toml").read_bytes()
    run("remove", "fixture-mr", "--no-refresh", "-y")
    assert not (root / "mods/fixture-mr.pw.toml").exists()
    assert (root / "mods/fixture-cf.pw.toml").is_file()
    assert (root / "index.toml").read_bytes() == before
    run("refresh")
    entries = tomllib.loads((root / "index.toml").read_text())["files"]
    assert [e["file"] for e in entries] == ["mods/fixture-cf.pw.toml"]
print("Offline provider metadata, deferred refresh, exact identity and batch removal passed")
