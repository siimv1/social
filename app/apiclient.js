export const apiRequest = async (endpoint, method = 'GET', body = null, parseJson = true) => {
  const headers = {
    'Content-Type': 'application/json',
  };

  try {
    const response = await fetch(`http://localhost:8080${endpoint}`, { 
      method,
      headers,
      body: body ? JSON.stringify(body) : null,
      credentials: 'include',  // Oluline küpsiste ja sessioonide jaoks
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorText}`);
    }

    if (parseJson) {
      return response.json();
    } else {
      return response;
    }
  } catch (error) {
    throw error;
  }
};
