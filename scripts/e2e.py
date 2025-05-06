import subprocess
import glob
import pytest
import sys

@pytest.mark.parametrize("case", glob.glob("../data/*"))
def test_case(tmp_path, case):
    cfg = f"{case}/config.json"
    events  = f"{case}/events"
    outlog = tmp_path / "out.log"

    result = subprocess.run(
        ["../../biathlon", "-config", cfg, "-events", events, "-out", str(outlog)],
        stdout=subprocess.PIPE,
        stderr=sys.stderr,
        text=True,
        check=True,
    )
    
    with open(f"{case}/expected.log") as f:
        expected_log = f.read().strip()
    actual_log = outlog.read_text().strip()
    assert actual_log == expected_log

    with open(f"{case}/expected.out") as f:
        expected_out = f.read().strip()
    actual_out = result.stdout.strip()
    assert actual_out == expected_out