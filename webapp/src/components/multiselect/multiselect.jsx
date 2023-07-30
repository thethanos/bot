import React, { useState } from "react";
import "./multiselect.css";

function Multiselect({ services, checked, handleCheck, handleClose, handleText }) {

    const [checkedLocal] = useState({ ...checked});
    const setChecked = (id)=> {
        checkedLocal[id] = checkedLocal[id]?!checkedLocal[id]:true;
      }

    let getText = ()=> {
        let text = "";
        for (let service of services) {
            if (checkedLocal[service.id]) {
                if (text.length === 0) {
                    text += service.name;
                } else {
                    text = text + ", " + service.name;
                }
            }
        }
        if (text.length === 0) {
            return "Выберите услуги";
        }
        return text;
    }

    return (
        <dialog className="modal-overlay" open>
            <div className="modal">
                <div className="modal-checkbox-container">
                    { services.map((service, index) => (
                            <div className="modal-checkbox-input-container" key={index + services.length * 2}>
                                <input 
                                    className="modal-checkbox" 
                                    key={index} 
                                    type="checkbox" 
                                    id={service.id}
                                    defaultChecked={checked[service.id]?checked[service.id]:false}
                                    onChange={()=>{
                                        setChecked(service.id);
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
                    <button onClick={()=>{handleClose(); handleText(getText()); handleCheck(checkedLocal)}}>Выбрать</button>
                    <button onClick={handleClose}>Отмена</button>
                </div>
            </div>
        </dialog>
    )
}

export default Multiselect;