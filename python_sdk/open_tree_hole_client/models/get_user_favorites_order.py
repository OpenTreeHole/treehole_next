from enum import Enum


class GetUserFavoritesOrder(str, Enum):
    HOLE_TIME_UPDATED = "hole_time_updated"
    ID = "id"
    TIME_CREATED = "time_created"

    def __str__(self) -> str:
        return str(self.value)
