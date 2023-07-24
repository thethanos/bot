import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import Gallery from './gallery';
import RegistrationForm from './registration_form';
import { BrowserRouter, Routes, Route} from 'react-router-dom'

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <BrowserRouter>
    <Routes>
      <Route path='/webapp/gallery' element={<Gallery />} />
      <Route path='/webapp/registration' element={<RegistrationForm />} />
    </Routes>
  </BrowserRouter>
);
