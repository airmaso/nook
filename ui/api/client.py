from abc import ABC, abstractmethod

import ui.config as config
from core.models import Leaderboard, User
from data.scraper import fetch_global_leaderboard, fetch_monthly_leaderboard


class NookApi(ABC):
    @abstractmethod
    def get_leaderboard(self, leaderboard: Leaderboard) -> list[User]:
        pass
    
    @abstractmethod
    def refresh(self, leaderboard: Leaderboard) -> float:
        pass


class LocalNookApi(NookApi):
    def get_leaderboard(self, leaderboard: Leaderboard) -> list[User]:
        if leaderboard == Leaderboard.GLOBAL:
            return fetch_global_leaderboard()
        elif leaderboard == Leaderboard.MONTHLY:
            return fetch_monthly_leaderboard()
        else:
            raise ValueError(
                config.ERR_API_INVALID_TYPE.format(type=leaderboard)
            )

    def refresh(self, leaderboard: Leaderboard) -> float:
        pass  # no-op since there's no cache to refresh


api = LocalNookApi()

if __name__ == "__main__":
    print(api.__class__)

    print(api.get_leaderboard(Leaderboard.MONTHLY))
    print(api.get_leaderboard(Leaderboard.GLOBAL))
    
    try:
        api.get_leaderboard("UNKNOWN")
    except ValueError as e:
        print(e)
