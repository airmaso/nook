from ui.api.client import api
from ui.app import Nook

if __name__ == "__main__":
    app = Nook(api=api)
    app.mainloop()
