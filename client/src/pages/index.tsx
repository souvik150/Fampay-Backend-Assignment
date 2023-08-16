import React, { useState, useEffect } from 'react';
import VideoList from "../../components/videoList";

export default function Home() {
  const [videos, setVideos] = useState([]);
  const [count, setCount] = useState();

  useEffect(() => {

      fetch('http://localhost/api/v1/video?page=1&limit=100')
          .then(response => response.json())
          .then(data => {
              setVideos(data.videos)
              setCount(data.totalResults);
              console.log(data.videos);
          })
          .catch(error => console.error('Error fetching videos:', error));
  }, []);

  return (
      <main className="flex min-h-screen flex-col items-center justify-between p-24">
          <h1>Total Count: {count}</h1>
          <VideoList videos={videos} />
      </main>
  );
}
