from enum import IntEnum


class ReportRange(IntEnum):
    RANGE_NOT_DEALT = 0
    RANGE_DEALT = 1
    RANGE_ALL = 2

    def __str__(self) -> str:
        return str(self.value)
