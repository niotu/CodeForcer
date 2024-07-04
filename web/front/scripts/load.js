export function exloadPage(page) {
    fetch('pages/' + page)
        .then(response => response.text())
        .then(html => {
            document.getElementById('content').innerHTML = html;
        })
        .catch(error => {
            console.error('Error loading page:', error);
        });
}