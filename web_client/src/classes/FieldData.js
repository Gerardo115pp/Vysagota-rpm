export const FieldStates = {
    NORMAL: 0,
    HAS_ERRORS: 1,
    READY: 2
}

export const verifyFormFields = (form_data) => {
    let field_value = null;
    let is_valid = true;
    for(let fd of form_data) {
        field_value = fd.getFieldValue();
        if (field_value === "") {
            
            if (fd.required) {
                fd.state = FieldStates.NORMAL;
                is_valid = false;
            } else {
                fd.state = FieldStates.READY;
            }
        } else if(!fd.isReady()) {

            fd.state = FieldStates.HAS_ERRORS;
            fd.error_message = "Invalid field";
            is_valid = false;
        } else {
            fd.state = FieldStates.READY;
        }
    }

    return is_valid;
}

class FieldData {
    constructor(field_id, validation_regex, name,type_name="text", required=true) {
        this.id = field_id;
        this.name = name;   
        this.regex = validation_regex;
        this.type = type_name;
        this.state = FieldStates.NORMAL;
        this.required = required;
        this.error_message = "hey! im error"; // link's adventure...
        this.placeholder = this.name;
    }

    clear = () => {
        this.getField().value = '';
    }

    getField = () => {
        return document.getElementById(this.id);
    }

    getFieldValue = () => {
        let field = this.getField();
        
        if (!field) {
            return "";
        }

        return this.getField().value;
    }

    isReady = () => {
        return this.regex.test(this.getFieldValue());
    }

    setPlaceholder = placeholder => {
        this.placeholder = placeholder;
    }
}

export default FieldData;