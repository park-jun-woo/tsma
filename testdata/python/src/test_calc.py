"""Content-based test fixture: references calc symbols by call, for D2 matching."""

from calc import classify, Calculator


def test_classify():
    assert classify(1) == "pos"
    assert classify(-1) == "nonpos"


def test_add():
    assert Calculator().add(1, 2) == 3
