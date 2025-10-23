from .exceptions import BlackboardAPIError, CourseNotFoundError

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .blackboard_client import BlackboardClient


class AnnouncementClient:
    def __init__(self, parent_client: "BlackboardClient"):
        self.parent = parent_client

    def get_announcements(self, course_id: str) -> list:
        pass

    def get_announcement(
        self,
    ) -> list:
        pass

    # TODO: Look at how I did the
    def update_announcement_text(self) -> list:
        pass


    #TODO: Test this
    def delete_announcements(self, course_id: str) -> list:
        _announcements = []
        url = self.parent.endpoints.get_announcements(course_id=course_id)
        response = self.parent.get(url=url)

        if response.status_code != 200:
            raise BlackboardAPIError(
                f"Failed to fetch announcements for {course_id}: {response.status_code}"
            )

        data = response.json().get("results", [])
        for a in data:
            if a.get("id"):
                url = self.parent.endpoints.delete_announcement(course_id=course_id)
                response = self.parent.delete(url=url)
                if response.status_code != 200:
                    raise BlackboardAPIError(
                        f"Failed to fetch announcements for {course_id}: {response.status_code}"
                    )
