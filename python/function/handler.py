from sanic.response import json


async def handle(request):
    # TODO implement
    return json({
        "statusCode": 200,
        "body": "Hello from OpenFaaS!"
    })
    