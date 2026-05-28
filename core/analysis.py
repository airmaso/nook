from core.models import User

def point_gap(u1: User, u2: User) -> int:
    return abs(u2.total_points - u1.total_points)

def puzzle_gap(u1: User, u2: User) -> tuple[float]:
    min_average = min(u1.average_points, u2.average_points)
    gap = point_gap(u1, u2)
    return (
        gap / 563,         # estimate based on max possible score
        gap / min_average  # estimate using lower ranked user's avg
    )

def time_gap(u1: User, u2: User) -> float:
    min_average_s_per_puzzle = min(u1.average_time, u2.average_time)
    avg_puzzles = puzzle_gap(u1, u2)[1]
    return avg_puzzles * min_average_s_per_puzzle

