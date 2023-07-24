import React from "react";
import "./multiselect.css"

function Multiselect({services, handleClose}) {
    return (
        <dialog className="modal-overlay" open>
            <div className="modal">
            {
                services.map((service, index)=>(
                    <div key={index} type="checkbox" id={service.id}>{service.name}</div>
                ))
            }
            <button onClick={handleClose}>Выбрать</button>
            <button onClick={handleClose}>Отмена</button>
            </div>
        </dialog>
    )
}

export default Multiselect;