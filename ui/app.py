import tkinter as tk
import ui.config as config


class Nook(tk.Tk):
    def __init__(self):
        super().__init__()

        self._setup()  # can't use _configure due to tk internals

    def _setup(self):
        self.title(config.APP_TITLE)
        self.geometry(f"{config.INIT_WIDTH}x{config.INIT_HEIGHT}")
        self.resizable(width=True, height=True)
