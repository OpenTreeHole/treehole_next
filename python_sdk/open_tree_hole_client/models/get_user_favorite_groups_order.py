from enum import Enum


class GetUserFavoriteGroupsOrder(str, Enum):
    ID = "id"
    TIME_CREATED = "time_created"
    TIME_UPDATED = "time_updated"

    def __str__(self) -> str:
        return str(self.value)
