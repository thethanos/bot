import React from "react";
import "./multiselect.css";

function Multiselect({ services, checked, handleCheck, handleClose }) {
    return (
        <dialog className="modal-overlay" open>
            <div className="modal">
                <div className="modal-checkbox-container">
                    { services.map((service, index) => (
                            <div className="modal-checkbox-input-container">
                                <input 
                                    className="modal-checkbox" 
                                    key={index} 
                                    type="checkbox" 
                                    id={service.id}
                                    defaultChecked={checked[service.id]?checked[service.id]:false}
                                    onChange={()=>{
                                        handleCheck(service.id);
                                    }}
                                />
                                <label
                                    className="modal-checkbox-label"
                                    key={index + services.length} 
                                    htmlFor={service.id}>
                                    {service.name}
                                    </label>
                            </div>
                        ))
                    }
                </div>
                <div className="modal-buttons-container">
                    <button onClick={handleClose}>Выбрать</button>
                    <button onClick={handleClose}>Отмена</button>
                </div>
            </div>
        </dialog>
    )
}

export default Multiselect;