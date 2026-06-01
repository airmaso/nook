import tkinter as tk

import ui.config as config
from ui.api.client import NookApi


class Nook(tk.Tk):
    api: NookApi

    def __init__(self, api: NookApi):
        super().__init__()

        self.api = api

        self._setup()  # can't use _configure due to tk internals

    def _setup(self):
        self.title(config.APP_TITLE)
        self.geometry(f"{config.INIT_WIDTH}x{config.INIT_HEIGHT}")
        self.resizable(width=True, height=True)
