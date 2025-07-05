   function toggleForm() {
        var form = document.getElementById("create-project-form");
        var button = document.getElementById("create-project-button");
        button.textContent = (form.classList.contains("show")) ? "Create new Project" : "Cancel";
        button.style.backgroundColor = (form.classList.contains("show")) ? "#4CAF50" : "#f44336";
        form.classList.toggle("show");
    }
    function deleteProject(projectId) {
        var projectCard = document.getElementById("project-" + projectId);
        if (confirm("Are you sure you want to delete this project?")) {
            fetch("/delete-project/" + projectId, {
                method: "DELETE"
            }).then(response => {
                if (response.ok) {
                    projectCard.remove();
                } else {
                    alert("Failed to delete project.");
                }
            });
        }
    }
    function loadProject(projectId) {
        window.location.href = "/projects/" + projectId;
    }
    async function showGitHubRepos() {
        const close_button = document.getElementById("close-github-repo-btn");
        close_button.style.display = "inline-block";
        close_button.style.backgroundColor = "#f44336";

        const listDiv = document.getElementById("github-repo-list");
        listDiv.innerHTML = "Loading...";
        listDiv.style.display = "block";
        const resp = await fetch("/user/repos");
        if (!resp.ok) {
            listDiv.innerHTML = "Failed to load repositories.";
            return;
        }
        const repos = await resp.json();
        listDiv.innerHTML = "<ul>" + repos.map(r =>
            `<li>
            <button class="repo-btn" onclick="openGithubModal('${r.full_name}', '${r.name}', '${r.description ? r.description.replace(/'/g, "\\'") : ''}')">
                ${r.full_name}
            </button>
        </li>`
        ).join("") + "</ul>";
    }
    function closeRepos() {
        const listDiv = document.getElementById("github-repo-list");
        listDiv.style.display = "none";
        const close_button = document.getElementById("close-github-repo-btn");
        close_button.style.display = "none";
    }
    function openGithubModal(fullName, name, desc) {
        document.getElementById("github-modal").style.display = "block";
        document.getElementById("modal-repo-name").textContent = fullName;
        document.getElementById("modal-repo-fullname").value = fullName;
        document.getElementById("modal-project-name").value = name;
        document.getElementById("modal-project-desc").value = desc || "";
    }
    function closeGithubModal() {
        document.getElementById("github-modal").style.display = "none";
    }
    function importGithubProject(event) {
        event.preventDefault();
        // Тут зробіть POST-запит на свій бекенд для створення проекту з GitHub
        // Наприклад:
        const form = document.getElementById("import-github-form");
        const data = new FormData(form);
        fetch("/import-github-project", {
            method: "POST",
            body: data
        }).then(resp => {
            if (resp.ok) {
                closeGithubModal();
                closeRepos();
                location.reload(); // або оновіть список проектів динамічно
            } else {
                alert("Не вдалося імпортувати проект");
            }
        });
        return false;
    }