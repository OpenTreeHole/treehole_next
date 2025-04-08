from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.common_http_error import CommonHttpError
from ...models.models_report import ModelsReport
from ...models.report_delete_model import ReportDeleteModel
from ...types import Response


def _get_kwargs(
    id: int,
    *,
    body: ReportDeleteModel,
) -> dict[str, Any]:
    headers: dict[str, Any] = {}

    _kwargs: dict[str, Any] = {
        "method": "delete",
        "url": f"/reports/{id}",
    }

    _body = body.to_dict()

    _kwargs["json"] = _body
    headers["Content-Type"] = "application/json"

    _kwargs["headers"] = headers
    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[Union[CommonHttpError, ModelsReport]]:
    if response.status_code == 200:
        response_200 = ModelsReport.from_dict(response.json())

        return response_200
    if response.status_code == 400:
        response_400 = CommonHttpError.from_dict(response.json())

        return response_400
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[Union[CommonHttpError, ModelsReport]]:
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
    body: ReportDeleteModel,
) -> Response[Union[CommonHttpError, ModelsReport]]:
    r"""Deal a report

     Mark a report as \"dealt\" and send notification to reporter

    Args:
        id (int):
        body (ReportDeleteModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CommonHttpError, ModelsReport]]
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
    body: ReportDeleteModel,
) -> Optional[Union[CommonHttpError, ModelsReport]]:
    r"""Deal a report

     Mark a report as \"dealt\" and send notification to reporter

    Args:
        id (int):
        body (ReportDeleteModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CommonHttpError, ModelsReport]
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
    body: ReportDeleteModel,
) -> Response[Union[CommonHttpError, ModelsReport]]:
    r"""Deal a report

     Mark a report as \"dealt\" and send notification to reporter

    Args:
        id (int):
        body (ReportDeleteModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[Union[CommonHttpError, ModelsReport]]
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
    body: ReportDeleteModel,
) -> Optional[Union[CommonHttpError, ModelsReport]]:
    r"""Deal a report

     Mark a report as \"dealt\" and send notification to reporter

    Args:
        id (int):
        body (ReportDeleteModel):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Union[CommonHttpError, ModelsReport]
    """

    return (
        await asyncio_detailed(
            id=id,
            client=client,
            body=body,
        )
    ).parsed
