from locust import HttpLocust, TaskSet, task, between


class UserBehavior(TaskSet):

    @task(2)
    def python(l):
        l.client.get("/api/python")

    @task(1)
    def ruby(l):
        l.client.post("/api/ruby")

    @task(3)
    def web(l):
        l.client.get("/")

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    wait_time = between(5.0, 9.0)
