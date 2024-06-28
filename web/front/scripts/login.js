function sendLoginRequest(event) {
    event.preventDefault(); 

    const key = document.getElementById('key').value;

    const secret = document.getElementById('secret').value;

    const handle = document.getElementById('handle').value;

    const password = document.getElementById('password').value;

    const data = {
		 key: key
		 secret: secret
        handle: handle,
        password: password
    };

fetch("http://localhost:8080/setAdmin", {
  method: "POST",
  body: JSON.stringify(data),
  headers: {"Content-type": "application/json"}
})
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log('Response from server:', data);
    })
    .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
    });

}

document.getElementById('loginForm').addEventListener('submit', sendLoginRequest);
