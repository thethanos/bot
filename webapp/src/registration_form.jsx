import React from "react";
import "./registration_form.css";
import { useEffect, useState } from "react";
import Multiselect from "./components/multiselect/multiselect";

function RegistrationForm() {

  const [checked, setChecks] = useState({});
  const [showMultiselect, setShowMultiselect] = useState(false)
  const [textMultiselect, setTextMultiselect] = useState("Выберите услугу");
  const [cities, setCities] = useState([]);
  const [serviceCategories, setServiceCategories] = useState([]);
  const [services, setServices] = useState([]);

  async function loadData(url, setter, dependency) {
    try {
      let response = await fetch(url);
      if (!response.ok) {
        console.error(`Error has occured during request GET ${url} ${response.status}`);
        return;
      }
      let json = await response.json();
      if (!json) {
        return;
      }

      setter(json);
      if (dependency) {
        dependency(json[0].id);
      }
    } catch (exception) {
      console.error(`Exception has been thrown during request GET ${url} ${exception}`);
    }
  }

  const onServiceCategoryLoad = (id) => {
    setChecks({});
    loadData(`https://bot-dev-domain.com:444/services?category_id=${id}`, setServices, null)
  }

  const onServiceCategoryChange = (event) => {
    let id = event.target.value;
    setChecks({});
    setTextMultiselect("Выберите услуги");
    loadData(`https://bot-dev-domain.com:444/services?category_id=${id}`, setServices, null)
  }

  useEffect(() => {
    loadData("https://bot-dev-domain.com:444/cities", setCities, null);
  }, [])

  useEffect(() => {
    loadData("https://bot-dev-domain.com:444/services/categories", setServiceCategories, onServiceCategoryLoad);
  }, [])

  const uploadFile = async (file, url) => {
    return new Promise(async (resolve, reject) => {
      try {
        let formData = new FormData();
        formData.append("file", file);

        let response = await fetch(url, {
          method: "POST",
          body: formData
        }
        );

        if (!response.ok) {
          console.error("Error has occured during request POST ", url, response.status);
          reject();
          return;
        }
      } catch (exception) {
        console.error("Exception has been thrown during file upload", exception);
        reject();
      }

      resolve();
    })
  }

  const onSubmit = async (event) => {
    event.preventDefault();

    let elements = event.target.elements;
    const nameInput = elements.name;
    const citySelect = elements.city;
    const serviceCategorySelect = elements.service_category;
    const imageInput = elements.images;
    const contactInput = elements.contact;
    const descriptionInput = elements.description;

    const services = [];
    for (let service in checked) {
      services.push(service);
    }

    const images = [];
    for (let image of imageInput.files) {
      images.push(image.name);
    }

    const object = {
      name: nameInput.value,
      city_id: citySelect.value,
      service_category_id: serviceCategorySelect.value,
      service_ids: services,
      contact: contactInput.value,
      description: descriptionInput.value,
      images: images
    };
    const body = JSON.stringify(object);

    try {
      let response = await fetch("https://bot-dev-domain.com:444/masters", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: body
      });

      let data = await response.json();
      if (!response.ok) {
        console.error("Error has occured during request POST https://bot-dev-domain.com:444/masters", response.status);
        return;
      }

      const promises = [];
      for (let file of imageInput.files) {
        promises.push(uploadFile(file, `https://bot-dev-domain.com:444/masters/images/${data.id}`));
      }
      Promise.all(promises);

      response = await fetch(`https://bot-dev-domain.com:444/masters/approve/${data.id}`, { method: 'POST' });
      if (!response.ok) {
        console.error(`Error has occured during request POST https://bot-dev-domain.com:444/masters/approve/${data.id}`, response.status);
        return;
      }

      nameInput.value = '';
      citySelect.selectedIndex = 0;
      serviceCategorySelect.selectedIndex = 0;
      setChecks({});
      setTextMultiselect("Выберите услуги");
      imageInput.value = '';
      contactInput.value = '';
      descriptionInput.value = '';
      window.alert("Регистрация прошла успешно!");

    } catch (exception) {
      console.error("Exception has been thrown during request POST https://bot-dev-domain.com:444/masters", exception);
    }
  }

  return (
    <div className="container">
      <form onSubmit={onSubmit}>
        <h1>Регистрация</h1>
        <p>Пожалуйста заполните анкету чтобы зарегистрироваться в системе в
          качестве мастера.</p>
        <hr />

        <label htmlFor="name"><b>Имя</b></label>
        <input type="text" placeholder="Введите свое имя" className="name" name="name" required />

        <label htmlFor="city"><b>Город</b></label>
        <select name="city" id="city" required>
          {cities && cities && cities.map((city, index) => (<option key={index} value={city.id}>{city.name}</ option>))}
        </select>

        <label htmlFor="service_category"><b>Категория услуг</b></label>
        <select className="service_category" name="service_category" required onChange={onServiceCategoryChange}>
          {serviceCategories && serviceCategories.map((category, index) => (<option key={index} value={category.id}>{category.name}</option>))}
        </select>

        <label htmlFor="services"><b>Услуга</b></label>
        <div className="services" onClick={() => { setShowMultiselect(true) }}>{textMultiselect}</div>
        {showMultiselect &&
          <Multiselect
            services={services}
            checked={checked}
            handleCheck={setChecks}
            handleClose={() => { setShowMultiselect(false) }}
            handleText={setTextMultiselect}
          />
        }
        <label htmlFor="images"><b>Фотографии</b></label>
        <input type="file" multiple name="images" id="images" accept="image/*" required />

        <label htmlFor="contact"><b>Контактные данные</b></label>
        <input type="text" placeholder="Введите номер телефона или ссылку на социальную сеть" name="contact" id="contact" required />

        <label htmlFor="description"><b>Коротко о себе</b></label>
        <input type="text" placeholder="Текст, который будет отображаться в вашем профиле" name="description" id="description" />

        <hr />
        <button type="submit" className="registerbtn">Зарегистрироваться</button>
      </form>
    </div>
  )
}

export default RegistrationForm;