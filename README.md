# Concourse Notifier Plus
A failing Concourse pipeline triggers a slack notification with info about the failure. The team can interact with that notification by typing /notifier-plus in the team slack channel. The app then starts a tmate session which hijacks the container of the failing job and gives the developer a URL that allows them to jump on that session and troubleshoot the job output.
