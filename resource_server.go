package ahvproviderplugin

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	st "github.com/ideadevice/terraform-ahv-provider-plugin/jsonstruct"
	"github.com/ideadevice/terraform-ahv-provider-plugin/requestutils"
	set "github.com/ideadevice/terraform-ahv-provider-plugin/setjsonfields"
	"io/ioutil"
	"log"
	"runtime/debug"
)

type specStruct struct {
	Name      string      `json:"name"`
	Resources interface{} `json:"resources"`
}

type metaStruct struct {
	OwnerReference interface{} `json:"owner_reference"`
	SpecVersion    int64       `json:"spec_version"`
	UUID           string      `json:"uuid"`
	Kind           string      `json:"kind"`
	Categories     interface{} `json:"categories"`
}

type vmStruct struct {
	Metadata metaStruct  `json:"metadata"`
	Status   interface{} `json:"status"`
	Spec     specStruct  `json:"spec"`
}

type vmList struct {
	APIVersion string      `json:"api_version"`
	MetaData   interface{} `json:"metadata"`
	Entities   []vmStruct  `json:"entities"`
}

func updateAddress(d *schema.ResourceData) error {
	return nil
}

// Function checks if there is an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// RecoverFunc can be used to recover from panics. name is the name of the caller
func RecoverFunc(name string) {
	if err := recover(); err != nil {
		log.Printf("Recovered from error %s, %s", err, name)
		log.Printf("Stack Trace: %s", debug.Stack())
		panic(err)
	}
}

// ID returns the id to be set
func (m *Machine) ID() string {
	return "ID-" + m.Spec.Name + "!!"
}

// DeleteMachine function deletes the vm using DELETE api call
func (c *MyClient) DeleteMachine(m *Machine) error {

	jsonStr := []byte(`{}`)
	url := c.Endpoint + "/list"
	method := "POST"
	jsonResponse := requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password)

	var uuid string
	var vmlist vmList
	err := json.Unmarshal(jsonResponse, &vmlist)
	check(err)

	for _, vm := range vmlist.Entities {
		if vm.Spec.Name == m.Spec.Name {
			uuid = vm.Metadata.UUID
		}
	}

	url = c.Endpoint + "/" + uuid
	method = "DELETE"
	requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password)
	return nil
}

// CreateMachine function creates the vm using POST api call
func (c *MyClient) CreateMachine(m *Machine, d *schema.ResourceData) error {

	var JSON st.JSONstruct

	Input, err := ioutil.ReadFile("json_template")
	check(err)
	InputPattern := []byte(Input)

	json.Unmarshal(InputPattern, &JSON)

	set.SetJSONFields(&JSON, d)
	JSON.Spec.Name = m.Spec.Name
	JSON.Metadata.Name = m.Metadata.Name

	jsonStr, err1 := json.Marshal(JSON)
	check(err1)

	method := "POST"
	requestutils.RequestHandler(c.Endpoint, method, jsonStr, c.Username, c.Password)
	return nil
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*MyClient)
	//specTemp := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{})
	//resourcesTemp := specTemp["resources"].(*schema.Set).List()[0].(map[string]interface{})

	machine := Machine{
		Spec: &st.SpecStruct{
			Name: d.Get("name").(string),
		},
		Metadata: &st.MetaDataStruct{
			Name: d.Get("name").(string),
		},
	}

	err := client.CreateMachine(&machine, d)
	if err != nil {
		return err
	}

	d.SetId(machine.ID())
	return nil

}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	// checking that address has changed or not
	if d.HasChange("address") {
		//Try updating the address
		if err := updateAddress(d); err != nil {
			return err
		}
		// After updating address
		d.SetPartial("address")
	}
	// If we were to return here, before disabling patial mode below, then only "address" field would be saved

	//Disabling partial state mode. This will cause terraform to save all fields again
	d.Partial(false)

	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*MyClient)
	//specTemp := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{})
	//resourcesTemp := specTemp["resources"].(*schema.Set).List()[0].(map[string]interface{})

	machine := Machine{
		Spec: &st.SpecStruct{
			Name: d.Get("name").(string),
		},
		Metadata: &st.MetaDataStruct{
			Name: d.Get("name").(string),
		},
	}

	err := client.DeleteMachine(&machine)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_vcpus_per_socket": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"num_sockets": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"memory_size_mb": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"power_state": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"api_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"kind": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"creation_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"spec_version": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"entity_version": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"categories": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"owner_reference": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
