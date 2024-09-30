export const apiRequest = async (endpoint, method = 'GET', body = null, parseJson = true) => {
  const token = localStorage.getItem('token');
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`, // Parandatud interpolatsioon siin
  };

  try {
    const response = await fetch(`http://localhost:8080${endpoint}`, { // Parandatud interpolatsioon siin
      method,
      headers,
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      // Handle HTTP errors
      const errorText = await response.text();
      throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorText}`); // Parandatud interpolatsioon siin
    }

    if (parseJson) {
      return response.json();
    } else {
      return response; // Return the raw response if parseJson is false
    }
  } catch (error) {
    throw error;
  }
};
