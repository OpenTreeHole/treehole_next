"""Contains all the data models used in inputs/outputs"""

from .common_error_detail_element import CommonErrorDetailElement
from .common_error_detail_element_value import CommonErrorDetailElementValue
from .common_http_error import CommonHttpError
from .division_create_model import DivisionCreateModel
from .division_delete_model import DivisionDeleteModel
from .division_modify_division_model import DivisionModifyDivisionModel
from .favourite_add_favorite_group_model import FavouriteAddFavoriteGroupModel
from .favourite_add_model import FavouriteAddModel
from .favourite_delete_favorite_group_model import FavouriteDeleteFavoriteGroupModel
from .favourite_delete_model import FavouriteDeleteModel
from .favourite_modify_favorite_group_model import FavouriteModifyFavoriteGroupModel
from .favourite_modify_model import FavouriteModifyModel
from .favourite_move_model import FavouriteMoveModel
from .favourite_response import FavouriteResponse
from .floor_ban_division import FloorBanDivision
from .floor_create_model import FloorCreateModel
from .floor_create_old_model import FloorCreateOldModel
from .floor_create_old_response import FloorCreateOldResponse
from .floor_delete_model import FloorDeleteModel
from .floor_modify_model import FloorModifyModel
from .floor_modify_model_like import FloorModifyModelLike
from .floor_modify_sensitive_floor_request import FloorModifySensitiveFloorRequest
from .floor_restore_model import FloorRestoreModel
from .floor_search_config_model import FloorSearchConfigModel
from .floor_sensitive_floor_response import FloorSensitiveFloorResponse
from .get_floors_sensitive_order_by import GetFloorsSensitiveOrderBy
from .get_holes_hole_id_floors_order_by import GetHolesHoleIdFloorsOrderBy
from .get_holes_hole_id_floors_sort import GetHolesHoleIdFloorsSort
from .get_reports_range import GetReportsRange
from .get_reports_sort import GetReportsSort
from .get_user_favorite_groups_order import GetUserFavoriteGroupsOrder
from .get_user_favorites_order import GetUserFavoritesOrder
from .get_users_me_floors_order_by import GetUsersMeFloorsOrderBy
from .get_users_me_floors_sort import GetUsersMeFloorsSort
from .gorm_deleted_at import GormDeletedAt
from .hole_create_model import HoleCreateModel
from .hole_create_old_model import HoleCreateOldModel
from .hole_create_old_response import HoleCreateOldResponse
from .hole_modify_model import HoleModifyModel
from .message_create_model import MessageCreateModel
from .models_division import ModelsDivision
from .models_favorite_group import ModelsFavoriteGroup
from .models_floor import ModelsFloor
from .models_floor_history import ModelsFloorHistory
from .models_hole import ModelsHole
from .models_hole_floors import ModelsHoleFloors
from .models_map import ModelsMap
from .models_message import ModelsMessage
from .models_message_data import ModelsMessageData
from .models_message_model import ModelsMessageModel
from .models_message_type import ModelsMessageType
from .models_punishment import ModelsPunishment
from .models_report import ModelsReport
from .models_tag import ModelsTag
from .models_user import ModelsUser
from .models_user_config import ModelsUserConfig
from .models_user_permission import ModelsUserPermission
from .models_user_permission_silent import ModelsUserPermissionSilent
from .penalty_forever_post_body import PenaltyForeverPostBody
from .penalty_post_body import PenaltyPostBody
from .report_add_model import ReportAddModel
from .report_ban_body import ReportBanBody
from .report_delete_model import ReportDeleteModel
from .report_range import ReportRange
from .subscription_add_model import SubscriptionAddModel
from .subscription_delete_model import SubscriptionDeleteModel
from .subscription_response import SubscriptionResponse
from .tag_create_model import TagCreateModel
from .tag_delete_model import TagDeleteModel
from .tag_modify_model import TagModifyModel
from .user_modify_model import UserModifyModel
from .user_user_config_model import UserUserConfigModel

__all__ = (
    "CommonErrorDetailElement",
    "CommonErrorDetailElementValue",
    "CommonHttpError",
    "DivisionCreateModel",
    "DivisionDeleteModel",
    "DivisionModifyDivisionModel",
    "FavouriteAddFavoriteGroupModel",
    "FavouriteAddModel",
    "FavouriteDeleteFavoriteGroupModel",
    "FavouriteDeleteModel",
    "FavouriteModifyFavoriteGroupModel",
    "FavouriteModifyModel",
    "FavouriteMoveModel",
    "FavouriteResponse",
    "FloorBanDivision",
    "FloorCreateModel",
    "FloorCreateOldModel",
    "FloorCreateOldResponse",
    "FloorDeleteModel",
    "FloorModifyModel",
    "FloorModifyModelLike",
    "FloorModifySensitiveFloorRequest",
    "FloorRestoreModel",
    "FloorSearchConfigModel",
    "FloorSensitiveFloorResponse",
    "GetFloorsSensitiveOrderBy",
    "GetHolesHoleIdFloorsOrderBy",
    "GetHolesHoleIdFloorsSort",
    "GetReportsRange",
    "GetReportsSort",
    "GetUserFavoriteGroupsOrder",
    "GetUserFavoritesOrder",
    "GetUsersMeFloorsOrderBy",
    "GetUsersMeFloorsSort",
    "GormDeletedAt",
    "HoleCreateModel",
    "HoleCreateOldModel",
    "HoleCreateOldResponse",
    "HoleModifyModel",
    "MessageCreateModel",
    "ModelsDivision",
    "ModelsFavoriteGroup",
    "ModelsFloor",
    "ModelsFloorHistory",
    "ModelsHole",
    "ModelsHoleFloors",
    "ModelsMap",
    "ModelsMessage",
    "ModelsMessageData",
    "ModelsMessageModel",
    "ModelsMessageType",
    "ModelsPunishment",
    "ModelsReport",
    "ModelsTag",
    "ModelsUser",
    "ModelsUserConfig",
    "ModelsUserPermission",
    "ModelsUserPermissionSilent",
    "PenaltyForeverPostBody",
    "PenaltyPostBody",
    "ReportAddModel",
    "ReportBanBody",
    "ReportDeleteModel",
    "ReportRange",
    "SubscriptionAddModel",
    "SubscriptionDeleteModel",
    "SubscriptionResponse",
    "TagCreateModel",
    "TagDeleteModel",
    "TagModifyModel",
    "UserModifyModel",
    "UserUserConfigModel",
)
