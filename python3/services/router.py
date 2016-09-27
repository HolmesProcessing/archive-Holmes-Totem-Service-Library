# imports for tornado
import tornado
from tornado import web, httpserver, ioloop
from python3.services.configuration import Metadata

# imports for info output
import os


class DummyHandler(tornado.web.RequestHandler):
    def get(self):
        #filename = self.get_argument("obj", strip=False)
        pass


def CreateInfoHandler(metadata):
    name        = str(metadata.name       ).replace("\n", "<br>")
    version     = str(metadata.version    ).replace("\n", "<br>")
    description = str(metadata.description).replace("\n", "<br>")
    copyright   = str(metadata.copyright  ).replace("\n", "<br>")
    license     = str(metadata.license    ).replace("\n", "<br>")
    class InfoHandler(tornado.web.RequestHandler):
        # Emits a string which describes the purpose of the analytics
        def get(self):
            info = """
<p>{name:s} - {version:s}</p>
<hr>
<p>{description:s}</p>
<hr>
<p>{license:s}</p>
<hr>
<p>{copyright:s}</p>
            """.strip().format(
                name        = name,
                version     = version,
                description = description,
                license     = license,
                copyright   = copyright
            )
            self.write(info)
    return InfoHandler


class Router(tornado.web.Application):
    def __init__(self, metadata, handlers):
        for key in ["description", "license"]:
            fpath = metadata.__getattribute__(key)
            if os.path.isfile(fpath):
                with open(fpath) as file:
                    metadata.__setattr__(key, file.read())

        handlers = [
            (r'/',          CreateInfoHandler(metadata)),
            (r'/analyze/',  handlers.get("analyze") or DummyHandler),
            (r'/feed/',     handlers.get("feed")    or DummyHandler),
            (r'/check/',    handlers.get("check")   or DummyHandler),
            (r'/results/',  handlers.get("results") or DummyHandler),
            (r'/status/',   handlers.get("status")  or DummyHandler),
        ]

        settings = dict(
            template_path=os.path.join(os.path.dirname(__file__), 'templates'),
            static_path=os.path.join(os.path.dirname(__file__), 'static'),
        )
        tornado.web.Application.__init__(self, handlers, **settings)
        self.engine = None

    def ListenAndServe(self, httpbinding):
        server = tornado.httpserver.HTTPServer(self)
        server.listen(httpbinding)
        tornado.ioloop.IOLoop.instance().start()
