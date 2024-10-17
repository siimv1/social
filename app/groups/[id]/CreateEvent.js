import React, { useState, useEffect } from 'react';
import { apiRequest } from '../../apiclient';

const CreateEvent = ({ groupId, userId, onEventCreated }) => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [eventTime, setEventTime] = useState('');
    const [events, setEvents] = useState([]);


    console.log('CreateEvent komponent imporditud ja käivitatud');


    const handleCreateEvent = async () => {
        try {
            const response = await apiRequest(`/groups/${groupId}/events`, 'POST', {
                creatorId: userId,
                title,
                description,
                eventTime,
            });

            console.log('Event created:', response);
            if (onEventCreated) {
                onEventCreated(response.event_id);
            }

           
            setTitle('');
            setDescription('');
            setEventTime('');

          
            loadEvents();
        } catch (error) {
            console.error('Failed to create event:', error);
        }
    };


    const loadEvents = async () => {
        try {
            const response = await apiRequest(`/groups/${groupId}/events`, 'GET');
            if (response && Array.isArray(response.events)) {
                setEvents(response.events);
            } else {
                setEvents([]);
            }
        } catch (error) {
            console.error('Failed to load events:', error);
        }
    };

    useEffect(() => {
        loadEvents();
    }, []);

    const handleRespond = async (eventId, responseOption) => {
        try {
            await apiRequest(`/events/${eventId}/respond`, 'POST', {
                userId,
                response: responseOption,
            });

            console.log('Response submitted');
            loadEvents();
        } catch (error) {
            console.error('Failed to respond to event:', error);
        }
    };

    return (
        <div className="create-event">
            <h3>Create Event</h3>
            <input
                type="text"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
            />
            <textarea
                placeholder="Description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
            />
            <input
                type="datetime-local"
                value={eventTime}
                onChange={(e) => setEventTime(e.target.value)}
                required
            />
            <button onClick={handleCreateEvent}>Create Event</button>

            <h3>Events</h3>
            {events.length > 0 ? (
                events.map((event) => (
                    <div key={event.id} className="event-item">
                        <h4>{event.title}</h4>
                        <p>{event.description}</p>
                        <p>When: {new Date(event.eventTime).toLocaleString()}</p>
                        <div className="event-options">
                            <button onClick={() => handleRespond(event.id, 'Going')}>Going</button>
                            <button onClick={() => handleRespond(event.id, 'Not Going')}>Not Going</button>
                        </div>
                    </div>
                ))
            ) : (
                <p>No events available.</p>
            )}
        </div>
    );
};

export default CreateEvent;