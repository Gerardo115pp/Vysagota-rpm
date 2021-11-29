import Main from './pages/Main/main.svelte';
import Login from './pages/Login/login.svelte';
import DoctorHub from './pages/DoctorHub/doctor-hub.svelte';
import PacientSummary from './pages/PacientsSummary/pacients-summary.svelte';
import PacientRegistration from './pages/PacientRegistration/pacient-registration.svelte';
import PacientHub from './pages/PacientsHub/pacients-hub.svelte';
import Mri from './pages/MriUploading/mri-uploading.svelte';

const routes = {
    '/': Main,
    '/login/:profile': Login,
    '/doctor': DoctorHub,
    '/pacients-summary': PacientSummary,
    '/pacient-registration': PacientRegistration,
    '/pacient': PacientHub,
    '/mri': Mri
}

export { routes };