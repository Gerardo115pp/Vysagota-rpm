<script>
    import Header from '../../components/Header.svelte';
    import Button from '../../components/Button.svelte';
    import vysagota_client from '../../classes/VysagotaClient';
    import Stack from '../../classes/stack';

    // <ui state variables/>
    let mris_images = [];
    let current_mri_image;

    // <operational variables/>
    let files_mount_point;
    let ia_service = vysagota_client.getEndpointData('ia_server');
    if (ia_service === undefined) {
        vysagota_client.getEndpoints(() => ia_service = vysagota_client.getEndpointData('ia_server') ); 
    }
        
    //<constants/>
    const mri_healthy_message = 'This MRI doenst show any signs of dementia.';
    const mri_dementia_message = 'We detected some anormalities in this MRI. Have in mind our system can make mistakes. please contact your assigned doctor for further information.';


    const processMRIs = (mris) => {
        const mris_stack = new Stack();
        mris_images = [];
        mris.forEach(mri => {
            mris_stack.push(mri);
        });
        requestDiagnosis(mris_stack);
    }

    const requestDiagnosis = (mris_stack) => {
        if (mris_stack.length > 0) {
            const mri = mris_stack.pop();
            const form_data = new FormData();
            form_data.append('mri', mri.blob, mri.name);
            const request = new Request(
                `http://${ia_service.host}:${ia_service.port}/diagnose`,
                {
                    method: 'POST',
                    body: form_data
                }
            );
            fetch(request)
                .then(response => response.json())
                .then(result => {
                    mri.positive = result.diagnosis === 1 ? false : true;
                    const new_mris =[...mris_images, mri];
                    mris_images = new_mris; // the pointer must change so svelte can re-render, otherwise if we just push the new mri to the stack, it wouldnt be re-rendered
                    if (!mris_stack.isEmpty()) {
                        requestDiagnosis(mris_stack);
                    }
                })
        }
    }

    const readFile = file => {
        return new Promise((onSuccess, onReject) => {
            let file_reader = new FileReader();
            file_reader.onload = () => onSuccess({src: file_reader.result, name: file.name, blob: file});
            file_reader.onerror = () => onReject(fr);
            file_reader.readAsDataURL(file);
        });
    }

    const savefiles = e => {
        const { target:file_input } = e;
        const promises = []
        for(let f of file_input.files) {
            if(/image.*/.test(f.type)) {
                promises.push(readFile(f))
            }
        }
        Promise.all(promises).then(files => {
            mris_images = files;
            processMRIs(files);
        });
    }

    $: current_mri_image = mris_images.length > 0 ? mris_images[0] : null;
</script>

<div id="mri-page" class="page">
    <Header return_to="/pacient"/>
    <main class="page-content-container">
        <aside class="page-content" id="mris-panel">
            <div class="upper-section" id="mris-list">
                {#each mris_images as mri, index}
                    <div on:click={() => current_mri_image = mris_images[index]} class={`mri-item ${!mri.positive ? 'clean' : ''}`}>
                        {mri.name}
                    </div>
                {:else}
                    <div class="mri-item no-mris">
                        No MRI images uploaded
                    </div>
                {/each}
                
            </div>
            <div class="hs" />
            <div class="bottom-section" id="controls-section">
                <Button 
                    onClick={() => files_mount_point.click()}
                    label="Upload MRI" 
                    isSquared={true} 
                    width="30%" 
                    padding="1.4vh 1.4vw" 
                />
                <div id="files-mounting-point">
                    <input
                        bind:this={files_mount_point}
                        on:change={savefiles}
                        type="file"
                        id="files-input"
                        multiple 
                    />
                </div>
            </div>
        </aside>
        <secion class="page-content" id="result-panel">
            <div class="upper-section" id="mri-display-container">
                <div id="mri-container">
                    {#if current_mri_image !== null}
                        <img src={current_mri_image.src} alt="svelte_idiota"/>
                    {:else}
                         <div id="no-mri-selected-msg">
                            No MRI image selected
                         </div>
                    {/if}
                </div>
            </div>
            <div class="hs" />
            <div class="bottom-section" id="status-message-container">
                {#if current_mri_image !== null}
                    <div class={current_mri_image.positive ? "negative" : "positive"} id="status-message">
                        {current_mri_image.positive ? mri_dementia_message : mri_healthy_message}
                    </div>
                {:else}
                    <div id="status-message">
                        ....
                    </div>
                {/if}
            </div>
        </secion>    
    </main>
</div>

<style>
    main.page-content-container {
        display: flex;
        width: 100%;
        flex-direction: row;
        justify-content: flex-start;
        align-items: stretch;
        user-select: none;
    }

    .page-content {
        display: flex;
        flex-direction: column;
    }

    .upper-section {
        height: 61vh;
    }

    .bottom-section {
        height: 30vh;
    }

    
    /*=============================================
    =            Results            =
    =============================================*/
    
    #result-panel {
        width: 80vw;
    }

    #mri-display-container{
        display: flex;
        justify-content: center;
        align-items: center;
    }

    #mri-container {
        display: flex;
        width: 60vw;
        background-color: var(--grey-color);
        height: calc(60vh - (2*5vh));
        justify-content: center;
        align-items: center;
        border-radius: 5px;
        border: 2px solid var(--theme-color);
    }

    #mri-container img {
        width: 35%;
    }

    #no-mri-selected-msg {
        font-size: 2rem;
        color: white;
    }

    #status-message-container {
        display: flex;
        justify-content: center;
        align-items: center;
    }


    #status-message {
        font-size: 2rem;
        text-align: center;
    }

    #status-message.positive {
        color: #488503;
    }

    #status-message.negative {
        color: var(--theme-color);
    }
    /*=============================================
    =            Mris             =
    =============================================*/
    
    #mris-panel {
        width: 20vw;
        border-right: 1px solid var(--theme-color);
    }

    .mri-item {
        cursor: pointer;
        font-family: var(--title-font);
        font-size: 1.4rem;
        padding: 2.4vh 1.4vw;
        text-align: center;
    }

    .mri-item.clean {
        color: #488503;
    }

    .mri-item:hover {
        background-color: var(--theme-color);
        color: white;
    }

    .no-mris {
        cursor: default;
        font-family: var(--body-font);
    }

    .no-mris:hover {
        background-color: transparent;
        color: var(--theme-color);
    }

    #controls-section {
        display: flex;
        justify-content: center;
        align-items: center;
    }
    
    #files-mounting-point {
        display: none;
    }
    
    

</style>