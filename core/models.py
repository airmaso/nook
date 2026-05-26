from enum import Enum
from dataclasses import dataclass
import datetime as dt

class Leaderboard(Enum):
    MONTHLY = "monthly"
    GLOBAL  = "global"

@dataclass
class User:
    uid:             int
    rank:            int
    username:        str
    last_active:     dt.datetime
    solved:          int
    attempted:       int
    acceptance_rate: float
    average_time:    float
    average_points:  float
    total_points:    int
    
    def __str__(self) -> str:
        rank = f"{self.rank:>3}"
        username = f"{'@{}'.format(self.username):<33}"
        total_points = f"{'{:,} {}'.format(self.total_points, 'pts'):<15}"

        return f"{rank} {username} {total_points}"

if __name__ == "__main__":
    airmaso = User(
        uid=0, rank=1, username="airmaso", last_active=dt.datetime.now(),
        solved=1919, attempted=1919, acceptance_rate=100.00,
        average_time=190.00, average_points=555.00, total_points=1_800_000
    )
    print(airmaso)
