const jd_address = JD_ADDRESS_ENV; // changed at build time by webpack

const  VysagotaApiStatus = {
    OK: 0,
    UNREACHABLE: 1,
    UNDER_FAILURE: 2,
    OVERLOADED: 3,
}

class VysagotaEndpoints {
    constructor() {
        console.log("VysagotaEndpoints");
        this.client_data = {
            endpoints: {},
            jd_status: VysagotaApiStatus.OK,
            jd_status_last_update: 0,
            vysagota_ready: false
        }
        this.getEndpoints();
    }

    getEndpoints = (callback) => {
        const request = new Request(`http://${jd_address}/endpoints`, {method: 'GET', mode: 'cors'});
        fetch(request)
            .then(response => response.json())
            .then(data => {
                this.client_data.endpoints = data;
                this.client_data.jd_status = VysagotaApiStatus.OK;
                this.client_data.jd_status_last_update = Date.now();
                this.client_data.vysagota_ready = true;
                if (callback) callback();
            }).catch(error => {
                console.log(error);
                this.client_data.jd_status = VysagotaApiStatus.UNREACHABLE;
            })
    }

    getEndpointData= (endpoint_name) => {
        return this.client_data.endpoints[endpoint_name];
    }

}

const vysagotaClient = new VysagotaEndpoints();
Object.freeze(vysagotaClient);
export default vysagotaClient;