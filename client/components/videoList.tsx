import React from 'react';

interface Video {
    id: number;
    title: string;
    description: string;
    publishTime: Date;
    thumbnailURL: string;
}

interface VideoListProps {
    videos: Video[];
}

const VideoList: React.FC<VideoListProps> = ({ videos }) => {
    return (
        <div className="video-list">
            {videos ? videos.map(video => (
                    <div key={video.id} className="video-item">
                        <img className="thumbnail" src={video.thumbnailURL} alt={video.title} />
                        <div className="video-info">
                            <h2 className="video-title">{video.title}</h2>
                            <p className="video-description">{video.description}</p>
                            <p className="publish-time">
                                Published at {new Date(video.publishTime).toLocaleString()}
                            </p>
                        </div>
                    </div>
                )) :
                <>
                    <h1>Loading...</h1>
                </>
            }
        </div>
    );
};

export default VideoList;
