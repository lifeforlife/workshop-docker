from sanic import Sanic
from sanic.log import logger
import os

import multiprocessing

import gunicorn.app.base

from function.handler import handle

app = Sanic(__name__)

async def event_handler(request):
    resp = await event_handler_path(request, '')
    return resp


async def event_handler_path(request, path):
    response_data = await handle(request)
    return response_data


app.add_route(event_handler, '/', methods=['GET', 'PUT', 'POST', 'PATCH', 'DELETE'])
app.add_route(event_handler_path, '/<path:path>', methods=['GET', 'PUT', 'POST', 'PATCH', 'DELETE'])


def number_of_workers():
    return (multiprocessing.cpu_count() * 2) + 1


class StandaloneApplication(gunicorn.app.base.BaseApplication):

    def __init__(self, app, options=None):
        self.options = options or {}
        self.application = app
        super(StandaloneApplication, self).__init__()

    def load_config(self):
        config = dict([(key, value) for key, value in self.options.items()
                       if key in self.cfg.settings and value is not None])
        for key, value in config.items():
            self.cfg.set(key.lower(), value)

    def load(self):
        return self.application


if __name__ == '__main__':
    options = {
        'bind': '%s:%s' % ('0.0.0.0', '5000'),
        'workers': number_of_workers(),
        'worker_class': 'sanic.worker.GunicornWorker',
    }
    StandaloneApplication(app, options).run()