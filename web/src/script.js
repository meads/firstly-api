
export async function login(username, password) {
    const response = await fetch('http://localhost:8080/login/', {
        mode:'cors',
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ username, password }),
    })

    if (response.ok) {
        const data = await response.json()
        sessionStorage.setItem("refreshToken", data.refreshToken)
        sessionStorage.setItem("accessToken", data.accessToken)
        sessionStorage.setItem("sessionId", data.sessionId)
        sessionStorage.setItem("userId", data.userId)

        return { success: true, data: data };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { success: false, data: `${response.status} : ${errorData}` };
}

export async function register(username, password) {    
    const response = await fetch('http://localhost:8080/register/', {
        mode:'cors',
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ username, password }),
    })
    
    if (response.ok) {
        const data = await response.json()
        sessionStorage.setItem("refreshToken", data.refreshToken)
        sessionStorage.setItem("accessToken", data.accessToken)
        sessionStorage.setItem("sessionId", data.sessionId)
        sessionStorage.setItem("userId", data.userId)

        return { success: true, data: data };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { success: false, data: `${response.status} : ${errorData}` };
}

export async function logout() {
    let sessionId = sessionStorage.getItem("sessionId")
    const response = await fetch(`http://localhost:8080/logout/${sessionId}`, {
        mode:'cors',
        method: 'POST',
    })

    if (response.ok) { 
        const data = await response.text()

        return { success: true, data: data };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { success: false, data: `${response.status} : ${errorData}` };
}

export async function createNote(title, content) {
    const userIdString = sessionStorage.getItem("userId");
    const userId = parseInt(userIdString, 10);
    const accessToken = sessionStorage.getItem("accessToken");
    const response = await fetch(`http://localhost:8080/users/${userIdString}/notes/`, {
        mode:'cors',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify({ title, content, userId }),
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data, tokenExpired: false };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { 
        success: false, 
        data: `${errorData}`,
        tokenExpired: tokenExpired(response.status, errorData),
    }
}

export async function updateNote(id, title, content) {
    const userIdString = sessionStorage.getItem("userId");
    const userId = parseInt(userIdString, 10);
    const accessToken = sessionStorage.getItem("accessToken");
    const response = await fetch(`http://localhost:8080/users/${userIdString}/notes/`, {
        mode:'cors',
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify({ id, title, content, userId }),
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data, tokenExpired: false };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { 
        success: false, 
        data: `${errorData}`,
        tokenExpired: tokenExpired(response.status, errorData),
    }
}

export async function deleteNote(id) {
    const userIdString = sessionStorage.getItem("userId");
    const accessToken = sessionStorage.getItem("accessToken");
    const response = await fetch(`http://localhost:8080/users/${userIdString}/notes/${id}`, {
        mode:'cors',
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { success: false, data: `${response.status} : ${errorData}` };
}

export async function getNotes() {
    const userIdString = sessionStorage.getItem("userId");
    const accessToken = sessionStorage.getItem("accessToken");
    const response = await fetch(`http://localhost:8080/users/${userIdString}/notes/`, {
        mode:'cors',
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data, tokenExpired: false };
    }
    
    let errorData = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await response.json();
        errorData = errorData.error;
    } else {
        errorData = await response.text();
    }
    
    return { 
        success: false, 
        data: `${errorData}`,
        tokenExpired: tokenExpired(response.status, errorData),
    }
}

function tokenExpired(status, data) {
    return status === 401 && data.includes("token expired")
}
