from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.common_http_error import CommonHttpError
from ...models.favourite_modify_favorite_group_model import FavouriteModifyFavoriteGroupModel
from ...models.models_favorite_group import ModelsFavoriteGroup
from ...types import Response


def _get_kwargs(
    *,
    body: FavouriteModifyFavoriteGroupModel,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "put",
        "url": "/user/favorite_groups",
    }

    _body = body.to_dict()

    _kwargs["json"] = _body
    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[CommonHttpError, list["ModelsFavoriteGroup"]]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = ModelsFavoriteGroup.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if response.status_code == 404:
        response_404 = CommonHttpError.from_dict(response.json())

        return response_404
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[CommonHttpError, list["ModelsFavoriteGroup"]]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    body: FavouriteModifyFavoriteGroupModel,
) -> Response[Union[CommonHttpError, list["ModelsFavoriteGroup"]]]:
    """Modify User's Favorite Group

    Args:
        body (FavouriteModifyFavoriteGroupModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CommonHttpError, list['ModelsFavoriteGroup']]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    body: FavouriteModifyFavoriteGroupModel,
) -> Optional[Union[CommonHttpError, list["ModelsFavoriteGroup"]]]:
    """Modify User's Favorite Group

    Args:
        body (FavouriteModifyFavoriteGroupModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CommonHttpError, list['ModelsFavoriteGroup']]
    """

    return sync_detailed(
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    body: FavouriteModifyFavoriteGroupModel,
) -> Response[Union[CommonHttpError, list["ModelsFavoriteGroup"]]]:
    """Modify User's Favorite Group

    Args:
        body (FavouriteModifyFavoriteGroupModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CommonHttpError, list['ModelsFavoriteGroup']]]
    """

    kwargs = _get_kwargs(
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    body: FavouriteModifyFavoriteGroupModel,
) -> Optional[Union[CommonHttpError, list["ModelsFavoriteGroup"]]]:
    """Modify User's Favorite Group

    Args:
        body (FavouriteModifyFavoriteGroupModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CommonHttpError, list['ModelsFavoriteGroup']]
    """

    return (
        await asyncio_detailed(
            client=client,
            body=body,
        )
    ).parsed
