from enum import Enum


class FloorModifyModelLike(str, Enum):
    ADD = "add"
    CANCEL = "cancel"

    def __str__(self) -> str:
        return str(self.value)
