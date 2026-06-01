import tkinter as tk
import tkinter.ttk as ttk

import ui.config as config
from core.models import Leaderboard
from ui.api.client import NookApi
from ui.base.frames import NookFrame
from ui.tabs.leaderboard import NookLeaderboardTab
from ui.theme.main import Theme, DEFAULT_THEME


class Nook(tk.Tk):
    api:         NookApi
    theme:       Theme
    notebook:    ttk.Notebook
    monthly_tab: NookFrame
    global_tab:  NookFrame

    def __init__(self, api: NookApi, theme: Theme = DEFAULT_THEME):
        super().__init__()

        self.api = api
        self.theme = theme

        self._setup()  # can't use _configure due to tk internals
        self._build()

    def _setup(self):
        self.title(config.APP_TITLE)
        self.geometry(f"{config.INIT_WIDTH}x{config.INIT_HEIGHT}")
        self.resizable(width=True, height=True)
        self.configure(bg=self.theme.palette.bg)

    def _build(self):
        self.theme.apply_ttk_styles(ttk.Style(self))

        self.notebook = ttk.Notebook(self, style="nook.TNotebook")
        self.notebook.pack(
            fill=tk.BOTH,
            expand=True,
            pady=self.theme.spacing.top_py
        )

        self.add_tab(NookLeaderboardTab(
            self.notebook,
            on_load=lambda: self.api.get_leaderboard(Leaderboard.MONTHLY)
        ), title=Leaderboard.MONTHLY.value.capitalize())

        self.add_tab(NookLeaderboardTab(
            self.notebook,
            on_load=lambda: self.api.get_leaderboard(Leaderboard.GLOBAL)
        ), title=Leaderboard.GLOBAL.value.capitalize())

    def add_tab(self, tab: NookLeaderboardTab, title: str):
        self.notebook.add(tab, text=title)
