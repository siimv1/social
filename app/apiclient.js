export const apiRequest = async (endpoint, method = 'GET', body = null) => {
    const token = localStorage.getItem('token'); // Võtame tokeni localStorage'st

    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`, // Lisame Authorization päisesse
    };

    const options = {
        method,
        headers,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(`http://localhost:8080${endpoint}`, options);

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Network response was not ok');
    }

    return response.json();
};
