from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.division_modify_division_model import DivisionModifyDivisionModel
from ...models.models_division import ModelsDivision
from ...models.models_message_model import ModelsMessageModel
from ...types import Response


def _get_kwargs(
    id: int,
    *,
    body: DivisionModifyDivisionModel,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "patch",
        "url": f"/divisions/{id}/_webvpn",
    }

    _body = body.to_dict()

    _kwargs["json"] = _body
    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[ModelsDivision, ModelsMessageModel]]:
    if response.status_code == 200:
        response_200 = ModelsDivision.from_dict(response.json())

        return response_200
    if response.status_code == 404:
        response_404 = ModelsMessageModel.from_dict(response.json())

        return response_404
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[ModelsDivision, ModelsMessageModel]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: DivisionModifyDivisionModel,
) -> Response[Union[ModelsDivision, ModelsMessageModel]]:
    """Modify A Division

    Args:
        id (int):
        body (DivisionModifyDivisionModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsDivision, ModelsMessageModel]]
    """

    kwargs = _get_kwargs(
        id=id,
        body=body,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: DivisionModifyDivisionModel,
) -> Optional[Union[ModelsDivision, ModelsMessageModel]]:
    """Modify A Division

    Args:
        id (int):
        body (DivisionModifyDivisionModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsDivision, ModelsMessageModel]
    """

    return sync_detailed(
        id=id,
        client=client,
        body=body,
    ).parsed


async def asyncio_detailed(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: DivisionModifyDivisionModel,
) -> Response[Union[ModelsDivision, ModelsMessageModel]]:
    """Modify A Division

    Args:
        id (int):
        body (DivisionModifyDivisionModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[ModelsDivision, ModelsMessageModel]]
    """

    kwargs = _get_kwargs(
        id=id,
        body=body,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    id: int,
    *,
    client: Union[AuthenticatedClient, Client],
    body: DivisionModifyDivisionModel,
) -> Optional[Union[ModelsDivision, ModelsMessageModel]]:
    """Modify A Division

    Args:
        id (int):
        body (DivisionModifyDivisionModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[ModelsDivision, ModelsMessageModel]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            body=body,
        )
    ).parsed
