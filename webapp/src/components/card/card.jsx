import React from "react";
import "./card.css"
import { useState } from "react";

function Card(properties) {

    let [image, setImage] = useState(0);
    let [activeDot, setActiveDot] = useState(0);
    const prevImage = () => {
        image--;
        if (image < 0) {
            image = properties.images.length - 1;
        }
        setImage(image);
        setActiveDot(image);
    }
    const nextImage = () => {
        image = (image + 1) % properties.images.length;
        setImage(image);
        setActiveDot(image);
    }

    return (
        <div className="master-container">
            <div className="master-image-container">
                <div className="master-images">
                    <img src={properties.images[image]} alt="nice pic" />
                </div>
                <button className="master-image-prev-btn" onClick={prevImage}></button>
                <button className="master-image-next-btn" onClick={nextImage}></button>
                <div className="dots">
                    {   properties.images.length > 1 &&
                        properties.images.map((_, index) => (
                            <span key={index} className={`dot ${index === activeDot ? 'dot-active' : 'dot-inactive'}`} />
                        )
                    )}
                </div>
            </div>
            <div className="master-bottom">
                <h3 className="master-name">{properties.name}</h3>
                <p className="master-description">{properties.description}</p>
                <button className="master-contact-btn">Написать мастеру</button>
            </div>
        </div>
    );
}

export default Card;

