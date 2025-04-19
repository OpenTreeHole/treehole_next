from enum import Enum


class GetHolesHoleIdFloorsOrderBy(str, Enum):
    ID = "id"
    LIKE = "like"

    def __str__(self) -> str:
        return str(self.value)
