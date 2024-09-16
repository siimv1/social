export const apiRequest = async (endpoint, method, body) => {
    const response = await fetch(`http://localhost:8080${endpoint}`, {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    });

    if (!response.ok) {
        throw new Error('Network response was not ok');
    }

    return response.json();
};

