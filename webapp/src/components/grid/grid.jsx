import React from "react";
import "./grid.css"
import Card from "../card/card"

function Grid(properties) {
    return (
        <div className="master-grid">
            {
                properties.cards.map((card) => (
                    <Card images={card.images} name={card.name} description={card.description} />
                )
            )}
        </div>
    )
}   

export default Grid;