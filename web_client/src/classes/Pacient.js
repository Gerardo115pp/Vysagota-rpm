
import vysagota_endpoints from "./VysagotaClient";

const validName = (name) => /^([A-z]\s?){1,3}$/.test(name);
const validUsername = username => /^[A-z_\-\.]{3,15}$/.test(username);
const validPassword = password => /^[A-z0-9\.#%]{8,}$/.test(password);
const validHash = hash => /^[A-z0-9]{64}$/.test(hash);
const validEmail = email => /^[A-z\d\._]{1,30}@([A-z0-9]{3,15}\.){2,15}$/.test(email);
const validPhone = phone => /^[0-9]{10}$/.test(phone);
const validAppointmentsCount = appointmentsCount => /^[\d]{1,}$/.test(appointmentsCount);
// we can test brithday with datetime which will cover more formats

export const getPacients = callback => {
    // TODO: clean this implementation, is to thight couple with getPacients
    const accounts_service = vysagota_endpoints.getEndpointData("accounts");
    if (accounts_service === undefined) {
        vysagota_endpoints.getEndpoints(() => getPacients(callback));
    } else {
        const request_headers = new Headers();
        request_headers.append("Cache-Control", "no-cache, no-store, must-revalidate");
        request_headers.append("Pragma", "no-cache");
    
        const request = new Request(`http://${accounts_service.host}:${accounts_service.port}/pacients`, {headers: request_headers, method: "GET"});
        return fetch(request)
                .then(promise => promise.json())
                .then(data => callback(data));
    }
   
}

class Pacient {
    constructor() {
        this.name = "generic-pacient";
        this.birthday = null;
        this.username = null;
        this.secret = null;
        this.email = null;
        this.tutor_email = null;
        this.phone = null;
        this.is_male = false;
        this.appointments = 0;
        this.last_activity = null;
        this.created_at = null;
    }

    create = () => {
        /*
            
            If pacient doenst exists, it creates it on the backend
            
        */
        const data = this.toJson();
        console.log(data);
        const form_data = new FormData();
        Object.entries(data).forEach((entrie) => {
            form_data.append(entrie[0], entrie[1]);
        });

        const accounts_service = vysagota_endpoints.getEndpointData("accounts");
        const request = new Request(`http://${accounts_service.host}:${accounts_service.port}/register`, {method: "POST", body: form_data});
        fetch(request);
    }

    toJson = () => {
        const data = {};
        data.name = this.name;
        data.username = this.username;
        data.password = this.secret;
        data.email = this.email;
        data.phone = this.phone;
        data.tutor_email = this.tutor_email;
        data.birthday = this.birthday.toISOString();
        data.appointments_count = this.appointments;
        if (this.last_activity === null) {
            data.last_activity = new Date().toISOString(); // if last_activity assume the pacient is been created
       } else {
           data.last_seen = this.last_activity.toISOString();
       }

        if (this.created_at === null) {
            data.created_at = new Date().toISOString();
        } else {
            data.created_at = this.created_at.toISOString();
        }

        data.gender = this.is_male ? 'M' : 'F';
        

        return data;
    }


    verifyAllFields = () => {
        let all_valid = true;
        let warning = "invalid pacient fields:\n"

        // NAME
        if (!validName(this.name)) {
            all_valid = false;
            warning += `name(${this.name})`;
        }


        // BRITHDAY
        try {
            new Date(this.birthday);
        } catch (error) {
            all_valid = false;
            warning += `brithday(${this.birthday})`;
        }

        // USERNAME
        if (!validUsername(this.username)) {
            all_valid = false;
            warning += `username(${this.username})`;
        }
        
        // SECRET
        if(!validPassword(this.secret)) {
            all_valid = false;
            warning += `secret(${this.secret})`;
        }

        // EMAIL
        if (!validEmail(this.email)) {
            all_valid = false;
            warning += `email(${this.email})`;
        }

        // TUTOR EMAIL
        if (!validEmail(this.tutor_email)) {
            all_valid = false;
            warning += `tutor_email(${this.tutor_email})`;
        }

        // PHONE
        if (!validPhone(this.phone)) {
            all_valid = false;
            warning += `phone(${this.phone})`;
        }

        // IS_MALE
        if (typeof this.is_male !== 'boolean') {
            all_valid = false;
            warning += `is_male(${this.is_male})`;
        }

        // APPOINTMENTS 
        if(typeof this.appointments !== 'number' || this.appointments < 0) {
            all_valid = false;
            warning += `appointments(${this.appointments})`;
        }

        // LAST ACTIVITY
        try {
            new Date(this.last_activity);
        } catch (error) {
            all_valid = false;
            warning += `last_activity(${this.last_activity})`;
        }
        
        if(!all_valid) {
            console.log(warning);
        }

        return all_valid;
    }

    exists = () => false;


}

export default Pacient;